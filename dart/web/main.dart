library blckur;

import 'package:blckur/routers/routers.dart' as routers;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/collections/events.dart' as evnts;
import 'package:blckur/components/components.dart' as components;
import 'package:blckur/decorators/decorators.dart' as decorators;
import 'package:blckur/formatters/formatters.dart' as formatters;
import 'package:blckur/models/models.dart' as models;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/application_factory.dart' as appfactory;
import 'package:angular/animate/module.dart' as animate;

void main() {
  logger.setup();
  var app = appfactory.applicationFactory()
    .addModule(new routers.RoutersMod())
    .addModule(new collections.CollectionsMod())
    .addModule(new components.ComponentsMod())
    .addModule(new decorators.DecoratorsMod())
    .addModule(new formatters.FormattersMod())
    .addModule(new models.ModelsMod())
    .addModule(new animate.AnimationModule())
    .run();

  app.get(evnts.Events).start();
}
