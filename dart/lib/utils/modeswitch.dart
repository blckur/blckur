library modeswitch;

import 'dart:async' as async;

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
        const Duration(milliseconds: 350), () {
          this._modeWait = null;
        });

      return;
    }

    this._modeWait = new async.Future.delayed(
      const Duration(milliseconds: 700), () {});

    new async.Timer(const Duration(milliseconds: 350), () {
      this.state[mode] = true;
      this.mode = mode;
      this._modeWait = null;
    });
  }

  int get mode {
    return this._mode;
  }
}
