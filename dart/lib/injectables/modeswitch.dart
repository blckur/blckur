library modeswitch_ing;

import 'package:angular/angular.dart' show Injectable;
import 'dart:async' as async;

@Injectable()
class ModeSwitch {
  int _mode;
  async.Future _modeWait;
  Map<int, bool> state;

  void set mode (int mode) {
    if (this._modeWait != null) {
      this._modeWait.then((_) {
        this._modeWait = null;
        this.mode = mode;
      });
      return;
    }

    if (this._mode == mode) {
      return;
    }
    this._mode = mode;

    this.state = {};

    if (mode == 0) {
      this._modeWait = new async.Future.delayed(
        const Duration(milliseconds: 400), () {
          this._modeWait = null;
        });

      return;
    }

    this._modeWait = new async.Future.delayed(
      const Duration(milliseconds: 800), () {});

    new async.Timer(const Duration(milliseconds: 400), () {
      this.state[mode] = true;
      this.mode = mode;
      this._modeWait = null;
    });
  }

  int get mode {
    return this._mode;
  }
}
