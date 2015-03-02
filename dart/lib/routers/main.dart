library main_rout;

import 'package:angular/angular.dart' as ng;

MainRout(ng.Router router, ng.RouteViewFactory views) {
  views.configure({
    'feed': ng.ngRoute(
      path: '/feed',
      view: 'packages/blckur/views/feed.html',
      defaultRoute: true
    ),
    'login': ng.ngRoute(
      path: '/login',
      view: 'packages/blckur/views/login.html'
    )
  });
}
