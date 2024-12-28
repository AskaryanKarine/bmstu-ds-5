# Настройка кубера для будущей себя

[Гайд по развертыванию кубера от timeweb](https://timeweb.cloud/tutorials/kubernetes/kak-ustanovit-i-nastroit-kubernetes-ubuntu?ysclid=m41ipzcbv4856716276)

[МеталЛБ для раздачи айпишников](https://metallb.universe.tf/installation/). Важно не забыть применить [этот манифест](ingress-nginx/ipadresspool.yaml) - раздача IPшников, обязательно должен быть указан пул 

Ингресс (самое больное тут) - для меньшей головной боли нужен именно ингресс от кубера **(Ingress-Nginx)**  
[Инструкция](https://sysadmin.info.pl/en/blog/how-to-install-nginx-as-ingress-controller-for-k3s/#install-nginx-as-ingress-controller-in-k3s)  
[Оф дока](https://kubernetes.github.io/ingress-nginx/deploy/?ref=blog.thenets.org#bare-metal-clusters)  
Важно не забыть про [этот манифест с настройкой LoadBalancer](ingress-nginx/ingress-controller-load-balancer.yaml)

Клоака:
[Оф дока, которая НЕ работает](https://www.keycloak.org/getting-started/getting-started-kube)  
Манифесты из оф доки лежат [тут](keycloak)  
[Понятный гайд по работе с клоакой](https://wkrzywiec.is-a.dev/posts/044_oauth2-keycloak/#adding-realm-client-roles-and-users)