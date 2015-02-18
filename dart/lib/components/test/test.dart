library test_comp;

import 'package:angular/angular.dart' show Component;

@Component(
  selector: 'test',
  templateUrl: 'packages/blckur/components/test/test.html',
  cssUrl: 'packages/blckur/components/test/test.css'
)
class TestComp {
  String inputValue;
  bool showToast;

  void onClick() {
    this.showToast = this.showToast != true;
  }
}
