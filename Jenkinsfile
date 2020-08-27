pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                // script {
                //     def browsers = ['chrome','firefox']
                //     for (int i = 0;i < browsers.size();++i){
                //         echo "Testing the ${browsers[i]} browser"
                //     }
                // }
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}
