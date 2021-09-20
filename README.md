# poc_go

Projeto simples para testar conceitos e tecnologias.

Para executar, instale as dependencias usando `make deps` e inicie os conteiners usando `docker-compose up -d`. Inicialmente, os serviços irão falhar, pois o etcd precisa ser populado com as configurações. Para popular o etcd com as configurações, execute o `make set-config` e inicie novamente os serviços parados com `docker-compose up -d`.
