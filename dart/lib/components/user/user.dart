library user_comp;

import 'package:angular/angular.dart' show Component;

@Component(
  selector: 'user',
  templateUrl: 'packages/blckur/components/user/user.html',
  cssUrl: 'packages/blckur/components/user/user.css'
)
class UserComp {
  bool settings;

  void onSettings() {
    settings = this.settings != true;
  }

  void onSave() {
    settings = false;
  }
}
