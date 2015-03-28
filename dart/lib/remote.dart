library remote;

import 'package:blckur/exceptions.dart';
import 'package:blckur/app.dart' as app;

import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;
import 'dart:html' as dom;

abstract class Remote {
  ng.Http http;
  ng.Router router;
  String url;
  String error;
  String errorMsg;
  int errorStatus;
  Function onImport;

  Remote() {
    this.http = app.root.get(ng.Http);
  }

  dynamic parse(dynamic data) {
    return data;
  }

  dynamic parseError(dynamic err) {
    var httpErr = new HttpError(err);

    this.error = httpErr.error;
    this.errorMsg = httpErr.errorMsg;
    this.errorStatus = httpErr.resp.status;

    if (this.errorStatus == 401) {
      dom.window.location.replace('#/login');
    }

    return httpErr;
  }

  void clearError() {
    this.error = null;
    this.errorMsg = null;
    this.errorStatus = null;
  }

  void import(dynamic responseData);

  void imported() {}

  async.Future fetch() {
    return this.http.get(this.url).then((response) {
      this.import(response.data);
      return response.data;
    }).catchError((err) {
      return new async.Future.error(this.parseError(err));
    }, test: (e) => e is ng.HttpResponse);
  }
}
