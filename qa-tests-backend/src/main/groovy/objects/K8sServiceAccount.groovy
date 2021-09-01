package objects

class K8sServiceAccount {
    String name
    String namespace
    Map<String, String> labels = [:]
    Map<String, String> annotations = [:]
    def automountToken
    def secrets = []
    String[] imagePullSecrets = []
}
