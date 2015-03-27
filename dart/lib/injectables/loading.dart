library loading;

import 'package:angular/angular.dart' show Injectable;
import 'dart:async' as async;

@Injectable()
class Loading {
  var _loading;
  bool state;

  bool set() {
    if (this._loading == true) {
      return false;
    }
    this._loading = true;

    new async.Timer(const Duration(milliseconds: 250), () {
      if (this._loading == true) {
        this.state = true;
      }
    });

    return true;
  }

  void clear() {
    this._loading = false;
    this.state = false;
  }
}
