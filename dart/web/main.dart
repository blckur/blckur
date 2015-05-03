library blckur;

import 'package:blckur/routers/routers.dart' as routers;
import 'package:blckur/collections/events.dart' as evnts;
import 'package:blckur/components/components.dart' as components;
import 'package:blckur/decorators/decorators.dart' as decorators;
import 'package:blckur/formatters/formatters.dart' as formatters;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/app.dart' as app;
import 'package:blckur/notification.dart' as notification;

import 'package:angular/angular.dart' as ng;
import 'package:angular/application_factory.dart' as appfactory;
import 'package:angular/animate/module.dart' as animate;
import 'dart:html' as dom;

void main() {
  logger.setup();
  app.root = appfactory.applicationFactory()
    .addModule(new routers.RoutersMod())
    .addModule(new components.ComponentsMod())
    .addModule(new decorators.DecoratorsMod())
    .addModule(new formatters.FormattersMod())
    .addModule(new animate.AnimationModule())
    .run();

  new evnts.Events(app.root.get(ng.RootScope),
    app.root.get(ng.Router)).start();

  dom.document.querySelector('.startup-background').remove();

  notification.init();
}
