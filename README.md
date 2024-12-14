# Secret Santa Backend

## Nomes dos Membros do Grupo

- *Joao Luiz Figueiredo Cerqueira*
- *Guilherme Soeiro de Carvalho Caporali*
- *Vítor Souza Fitzherbert*

## Explicação do Sistema

Este sistema consiste em uma API para uma aplicação de amigo secreto chamada *Secret Santa*. A API permite gerenciar grupos e participantes, realizar sorteios automáticos para definir os pares do amigo secreto e consultar informações relacionadas aos grupos e participantes. As principais funcionalidades incluem criar, atualizar, deletar grupos, adicionar participantes, gerar os matches automaticamente e consultar informações de forma segura e organizada.

### Principais Rotas Implementadas

- *POST /group* - Cria um novo grupo.
- *GET /group/:id* - Obtém detalhes de um grupo específico.
- *PUT /group/:id* - Atualiza um grupo existente.
- *DELETE /group/:id* - Remove um grupo existente.
- *POST /group/:id/add-participant* - Adiciona um participante a um grupo.
- *POST /group/:id/match-participants* - Realiza o sorteio dos participantes do grupo.
- *GET /group/:id/my-match* - Consulta o par atribuído a um participante.
- *GET /group* - Obtém todos os grupos cadastrados.

A estrutura de rotas foi configurada utilizando o framework *Gin*, permitindo uma organização clara e eficiente das requisições HTTP.

## Explicação das Tecnologias Utilizadas

- *Go Lang:* Linguagem principal usada para desenvolver a API devido à sua eficiência e robustez.

- *Gin Framework:* Framework usado para construir APIs RESTful rápidas e minimalistas.

- *Uber Dependency Injection:* Biblioteca usada para injeção de dependências, facilitando a manutenção e escalabilidade do sistema.

- *MongoDB:* Banco de dados NoSQL utilizado para armazenar dados de forma flexível e escalável.

- *Dockerfile e Docker Compose:* Utilizados para criação de imagens Docker e configuração do ambiente de desenvolvimento e produção.

- *Go Mock:* Utilizado para criação de mocks e simulação de dependências em testes unitários.

- *Swagger com Swaggo:* Usado para gerar automaticamente a documentação da API.

- *"go.mongodb.org/mongo-driver/mongo/integration/mtest":* Biblioteca para mockar operações do MongoDB em testes unitários.

- *Validation com "github.com/invopop/validation":* Biblioteca para validação de DTOs e estruturas de dados recebidas pela API.

- *Erros Customizados:* Implementação própria para gerenciamento de erros específicos da aplicação, facilitando o tratamento e a depuração de problemas.

- *CI para Testes Unitários:* Integração contínua configurada para executar automaticamente os testes unitários a cada atualização no repositório, garantindo a qualidade e estabilidade do sistema.

Este conjunto de tecnologias foi escolhido para garantir uma API robusta, escalável e segura, atendendo às melhores práticas de desenvolvimento de software moderno.
