library decorators;

import 'package:blckur/decorators/delay_click.dart' as delay_click;
import 'package:blckur/decorators/noselect.dart' as noselect;
import 'package:blckur/decorators/polymer.dart' as polymer;
import 'package:blckur/decorators/touch_click.dart' as touch_click;

import 'package:angular/angular.dart' as ng;

class DecoratorsMod extends ng.Module {
  DecoratorsMod() {
    this.bind(delay_click.DelayClickDec);
    this.bind(noselect.NoSelectDec);
    this.bind(polymer.PolymerDec);
    this.bind(touch_click.TouchClickDec);
  }
}
