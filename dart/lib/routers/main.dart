library main_rout;

import 'package:angular/angular.dart' as ng;

MainRout(ng.Router router, ng.RouteViewFactory views) {
  views.configure({
    'root': ng.ngRoute(
      path: '/feed',
      view: 'packages/blckur/views/root.html',
      defaultRoute: true
    ),
    'login': ng.ngRoute(
      path: '/login',
      view: 'packages/blckur/views/login.html'
    )
  });
}
