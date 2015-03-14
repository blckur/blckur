library loading;

import 'dart:async' as async;

class Loading {
  var _loading;
  bool loading;

  bool setLoading() {
    if (this._loading == true) {
      return false;
    }
    this._loading = true;

    new async.Timer(const Duration(milliseconds: 250), () {
      if (this._loading == true) {
        this.loading = true;
      }
    });

    return true;
  }

  void clearLoading() {
    this._loading = false;
    this.loading = false;
  }
}
