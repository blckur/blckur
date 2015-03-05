library remote;

import 'package:blckur/exceptions.dart';
import 'package:blckur/utils/utils.dart' as utils;

import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;
import 'dart:html' as dom;

abstract class Remote {
  String _loadCheckId;
  String _fetchCheckId;
  ng.Http http;
  ng.Router router;
  String url;
  String error;
  String errorMsg;
  int errorStatus;
  bool loading;
  bool loadingLong;
  Function onImport;

  Remote(this.http);

  String setLoading() {
    var loadCheckId = utils.uuid();
    this._loadCheckId = loadCheckId;
    this.loading = true;

    new async.Future.delayed(
      const Duration(milliseconds: 250), () {
        if (this._loadCheckId == loadCheckId) {
          this.loadingLong = true;
        }
      });

    return loadCheckId;
  }

  void clearLoading(String loadingId) {
    if (this._loadCheckId == loadingId) {
      this._loadCheckId = null;
      this.loadingLong = false;
      this.loading = false;
    }
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
    var fetchCheckId = utils.uuid();
    this._fetchCheckId = fetchCheckId;
    var loadId = this.setLoading();

    return this.http.get(this.url).then((response) {
      this.clearLoading(loadId);
      this.import(response.data);
      return response.data;
    }).catchError((err) {
      this.clearLoading(loadId);
      return new async.Future.error(this.parseError(err));
    }, test: (e) => e is ng.HttpResponse);
  }
}
