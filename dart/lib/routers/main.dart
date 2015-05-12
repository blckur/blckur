library main_rout;

import 'package:blckur/alert.dart' as alert;
import 'package:blckur/app.dart' as app;

import 'package:angular/angular.dart' as ng;
import 'dart:js' as js;

MainRout(ng.Router router, ng.RouteViewFactory views) {
  if (js.context['defaultPage'] == 'feed') {
    router.gotoUrl('/feed');
  }

  views.configure({
    'root': ng.ngRoute(
      defaultRoute: true,
      enter: (_) {
        var page = js.context['defaultPage'];

        if (page != null && page != '') {
          router.gotoUrl('/' + page);
        } else {
          router.gotoUrl('/feed');
        }
      }
    ),
    'feed': ng.ngRoute(
      path: '/feed',
      view: 'packages/blckur/views/feed.html',
      enter: (_) {
        if (app.events != null) {
          app.events.start();
        }
        alert.clear();
      },
      leave: (_) {
        if (app.events != null) {
          app.events.stop();
        }
      }
    ),
    'login': ng.ngRoute(
      path: '/login',
      view: 'packages/blckur/views/auth.html',
      enter: (_) {
        alert.clear();
      }
    ),
    'signup': ng.ngRoute(
      path: '/signup',
      view: 'packages/blckur/views/auth.html',
      enter: (_) {
        alert.clear();
      }
    ),
    'forgot': ng.ngRoute(
      path: '/forgot',
      view: 'packages/blckur/views/auth.html',
      enter: (_) {
        alert.clear();
      }
    ),
    'reset': ng.ngRoute(
      path: '/reset',
      view: 'packages/blckur/views/auth.html',
      enter: (_) {
        alert.clear();
      }
    )
  });
}
