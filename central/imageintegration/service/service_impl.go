package service

import (
	"context"
	"sort"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	clusterDatastore "github.com/stackrox/rox/central/cluster/datastore"
	"github.com/stackrox/rox/central/imageintegration/datastore"
	"github.com/stackrox/rox/central/reprocessor"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/or"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/images/integration"
	"github.com/stackrox/rox/pkg/registries"
	"github.com/stackrox/rox/pkg/scanners"
	"github.com/stackrox/rox/pkg/secrets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		or.SensorOrAuthorizer(user.With(permissions.View(resources.ImageIntegration))): {
			"/v1.ImageIntegrationService/GetImageIntegration",
			"/v1.ImageIntegrationService/GetImageIntegrations",
		},
		user.With(permissions.Modify(resources.ImageIntegration)): {
			"/v1.ImageIntegrationService/PostImageIntegration",
			"/v1.ImageIntegrationService/PutImageIntegration",
			"/v1.ImageIntegrationService/TestImageIntegration",
			"/v1.ImageIntegrationService/DeleteImageIntegration",
			"/v1.ImageIntegrationService/UpdateImageIntegration",
			"/v1.ImageIntegrationService/TestUpdatedImageIntegration",
		},
	})
)

// ImageIntegrationService is the struct that manages the ImageIntegration API
type serviceImpl struct {
	registryFactory registries.Factory
	scannerFactory  scanners.Factory
	toNotify        integration.ToNotify

	datastore        datastore.DataStore
	clusterDatastore clusterDatastore.DataStore
	reprocessorLoop  reprocessor.Loop
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterImageIntegrationServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterImageIntegrationServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

func scrubImageIntegration(i *storage.ImageIntegration) {
	secrets.ScrubSecretsFromStructWithReplacement(i, secrets.ScrubReplacementStr)
}

// GetImageIntegration returns the image integration given its ID.
func (s *serviceImpl) GetImageIntegration(ctx context.Context, request *v1.ResourceByID) (*storage.ImageIntegration, error) {
	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "image integration id must be provided")
	}
	integration, exists, err := s.datastore.GetImageIntegration(ctx, request.GetId())
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "image integration %s not found", request.GetId())
	}
	scrubImageIntegration(integration)
	return integration, nil
}

// GetImageIntegrations returns all image integrations.
func (s *serviceImpl) GetImageIntegrations(ctx context.Context, request *v1.GetImageIntegrationsRequest) (*v1.GetImageIntegrationsResponse, error) {
	integrations, err := s.datastore.GetImageIntegrations(ctx, request)
	if err != nil {
		return nil, err
	}

	identity := authn.IdentityFromContext(ctx)
	if identity != nil {
		svc := identity.Service()
		if svc != nil && svc.GetType() == storage.ServiceType_SENSOR_SERVICE {
			return &v1.GetImageIntegrationsResponse{Integrations: integrations}, nil
		}
	}

	// Remove secrets for other API accessors.
	for _, i := range integrations {
		scrubImageIntegration(i)
	}
	return &v1.GetImageIntegrationsResponse{Integrations: integrations}, nil
}

func sortCategories(categories []storage.ImageIntegrationCategory) {
	sort.SliceStable(categories, func(i, j int) bool {
		return int32(categories[i]) < int32(categories[j])
	})
}

func (s *serviceImpl) validateTestAndNormalize(ctx context.Context, request *storage.ImageIntegration) error {
	if err := s.validateIntegration(ctx, request); err != nil {
		return err
	}

	if !request.GetSkipTestIntegration() {
		if err := s.testImageIntegration(request); err != nil {
			return err
		}
	}

	sortCategories(request.Categories)
	return nil
}

// PutImageIntegration modifies a given image integration, without stored credential reconciliation
func (s *serviceImpl) PutImageIntegration(ctx context.Context, imageIntegration *storage.ImageIntegration) (*v1.Empty, error) {
	return s.UpdateImageIntegration(ctx, &v1.UpdateImageIntegrationRequest{Config: imageIntegration, UpdatePassword: true})
}

// PostImageIntegration creates a image integration.
func (s *serviceImpl) PostImageIntegration(ctx context.Context, request *storage.ImageIntegration) (*storage.ImageIntegration, error) {
	if request.GetId() != "" {
		return nil, status.Error(codes.InvalidArgument, "id field should be empty when posting a new image integration")
	}

	if err := s.validateTestAndNormalize(ctx, request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.datastore.AddImageIntegration(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	request.Id = id

	if err := s.toNotify.NotifyUpdated(request); err != nil {
		_ = s.datastore.RemoveImageIntegration(ctx, request.GetId())
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.reprocessorLoop.ShortCircuit()
	return request, nil
}

// DeleteImageIntegration removes a image integration given its ID.
func (s *serviceImpl) DeleteImageIntegration(ctx context.Context, request *v1.ResourceByID) (*v1.Empty, error) {
	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "image integration id must be provided")
	}
	if err := s.datastore.RemoveImageIntegration(ctx, request.GetId()); err != nil {
		return nil, err
	}
	if err := s.toNotify.NotifyRemoved(request.GetId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.Empty{}, nil
}

// UpdateImageIntegration modifies a given image integration, with optional stored credential reconciliation.
func (s *serviceImpl) UpdateImageIntegration(ctx context.Context, request *v1.UpdateImageIntegrationRequest) (*v1.Empty, error) {
	if err := s.validateIntegration(ctx, request.GetConfig()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.reconcileUpdateImageIntegrationRequest(ctx, request); err != nil {
		return nil, err
	}
	if err := s.validateTestAndNormalize(ctx, request.GetConfig()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.datastore.UpdateImageIntegration(ctx, request.GetConfig()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := s.toNotify.NotifyUpdated(request.GetConfig()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.reprocessorLoop.ShortCircuit()
	return &v1.Empty{}, nil
}

// TestImageIntegration checks if the given image integration is correctly configured, without using stored credential reconciliation.
func (s *serviceImpl) TestImageIntegration(ctx context.Context, imageIntegration *storage.ImageIntegration) (*v1.Empty, error) {
	return s.TestUpdatedImageIntegration(ctx, &v1.UpdateImageIntegrationRequest{Config: imageIntegration, UpdatePassword: true})
}

// TestImageIntegration checks if the given image integration is correctly configured, with optional stored credential reconciliation.
func (s *serviceImpl) TestUpdatedImageIntegration(ctx context.Context, request *v1.UpdateImageIntegrationRequest) (*v1.Empty, error) {
	if err := s.validateIntegration(ctx, request.GetConfig()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.reconcileUpdateImageIntegrationRequest(ctx, request); err != nil {
		return nil, err
	}
	if err := s.testImageIntegration(request.GetConfig()); err != nil {
		return nil, err
	}
	return &v1.Empty{}, nil
}

func (s *serviceImpl) testImageIntegration(request *storage.ImageIntegration) error {
	for _, category := range request.GetCategories() {
		if category == storage.ImageIntegrationCategory_REGISTRY {
			if err := s.testRegistryIntegration(request); err != nil {
				return status.Error(codes.InvalidArgument, errors.Wrap(err, "registry integration").Error())
			}
		}
		if category == storage.ImageIntegrationCategory_SCANNER {
			if err := s.testScannerIntegration(request); err != nil {
				return status.Error(codes.InvalidArgument, errors.Wrap(err, "scanner integration").Error())
			}
		}
	}
	return nil
}

func (s *serviceImpl) testRegistryIntegration(integration *storage.ImageIntegration) error {
	registry, err := s.registryFactory.CreateRegistry(integration)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err := registry.Test(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (s *serviceImpl) testScannerIntegration(integration *storage.ImageIntegration) error {
	scanner, err := s.scannerFactory.CreateScanner(integration)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err := scanner.Test(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (s *serviceImpl) validateIntegration(ctx context.Context, request *storage.ImageIntegration) error {
	if request == nil {
		return errors.New("empty integration")
	}
	errorList := errorhelpers.NewErrorList("Validation")
	if len(request.GetCategories()) == 0 {
		errorList.AddStrings("integrations require a category")
	}

	clustersRequested := request.GetClusters()
	existingClusters, err := s.clusterDatastore.GetClusters(ctx)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for _, req := range clustersRequested {
		if !s.clusterExists(req, existingClusters) {
			errorList.AddStringf("cluster %s does not exist", req)
		}
	}

	// Validate if there is a name. If there isn't, then skip the DB name check by returning the accumulated errors
	if request.GetName() == "" {
		errorList.AddString("name for integration is required")
		return errorList.ToError()
	}

	integrations, err := s.datastore.GetImageIntegrations(ctx, &v1.GetImageIntegrationsRequest{Name: request.GetName()})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if len(integrations) != 0 && request.GetId() != integrations[0].GetId() {
		errorList.AddStringf("integration with name %q already exists", request.GetName())
	}
	return errorList.ToError()
}

func (s *serviceImpl) clusterExists(name string, clusters []*storage.Cluster) bool {
	for _, c := range clusters {
		if name == c.GetName() {
			return true
		}
	}
	return false
}

func (s *serviceImpl) reconcileUpdateImageIntegrationRequest(ctx context.Context, updateRequest *v1.UpdateImageIntegrationRequest) error {
	if updateRequest.GetUpdatePassword() {
		return nil
	}
	if updateRequest.GetConfig() == nil {
		return status.Error(codes.InvalidArgument, "request is missing image integration config")
	}
	if updateRequest.GetConfig().GetId() == "" {
		return status.Error(codes.InvalidArgument, "id required for stored credential reconciliation")
	}
	integration, exists, err := s.datastore.GetImageIntegration(ctx, updateRequest.GetConfig().GetId())
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return status.Errorf(codes.NotFound, "image integration %s not found", updateRequest.GetConfig().GetId())
	}
	if err := s.reconcileImageIntegrationWithExisting(updateRequest.GetConfig(), integration); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (s *serviceImpl) reconcileImageIntegrationWithExisting(updatedConfig, storedConfig *storage.ImageIntegration) error {
	if updatedConfig.GetIntegrationConfig() == nil {
		return errors.New("the request doesn't have a valid integration config type")
	}
	return secrets.ReconcileScrubbedStructWithExisting(updatedConfig, storedConfig)
}
