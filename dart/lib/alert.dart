library alert;

import 'dart:collection' as collection;
import 'dart:async' as async;

collection.Queue<Alert> alerts = new collection.Queue();

class Alert {
  Function retryCallback;
  String text;

  Alert(this.text, [this.retryCallback]) {
    this.opened = true;

    alerts.add(this);

    while (alerts.length > 3) {
      alerts.removeFirst();
    }
  }

  bool _opened;
  bool get opened {
    return this._opened;
  }
  void set opened(bool val) {
    this._opened = val;
    if (!val) {
      new async.Timer(const Duration(milliseconds: 100), () {
        alerts.remove(this);
      });
    }
  }

  int get duration {
    if (this.retryCallback != null) {
      return 3600000;
    }
    return 4000;
  }
}

void clear() {
  alerts.clear();
}
