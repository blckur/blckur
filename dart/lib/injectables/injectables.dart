library injectables;

import 'package:blckur/injectables/loading.dart' as loading;
export 'package:blckur/injectables/loading.dart' show Loading;
import 'package:blckur/injectables/modeswitch.dart' as modeswitch;
export 'package:blckur/injectables/modeswitch.dart' show ModeSwitch;

import 'package:angular/angular.dart' as ng;

class InjectablesMod extends ng.Module {
  InjectablesMod() {
    this.bind(loading.Loading);
    this.bind(modeswitch.ModeSwitch);
  }
}
