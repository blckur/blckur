library model;

import 'package:blckur/remote.dart' as remote;

import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;

class Invalid extends Error {
  String type;
  String message;

  Invalid(this.type, this.message);

  toString() => this.message;
}

abstract class Model extends remote.Remote {
  dynamic id;
  Function onLinkClear;

  Model() : super() {
    this.init();
  }

  Model newModel() {
    throw new UnimplementedError('Model new not implemented.');
  }

  Map<String, Function> get getters {
    throw new UnimplementedError('Getter map not implemented.');
  }

  Map<String, Function> get setters {
    throw new UnimplementedError('Setter map not implemented.');
  }

  Map<String, Function> get validators {
    return {};
  }

  void validate(String name) {
    var validator = this.validators[name];

    if (validator != null) {
      validator(this.getters[name]());
    }
  }

  Model clone() {
    var mdl = this.newModel();
    var getters = this.getters;
    var setters = mdl.setters;

    setters.forEach((name, setter) {
      setter(getters[name]());
    });

    return mdl;
  }

  void import(dynamic responseData) {
    var data = this.parse(responseData);
    var setters = this.setters;

    if (data != null && data != '') {
      data.forEach((key, value) {
        var setter = setters[key];

        if (setter != null) {
          setter(value);
        }
      });
    }

    this.imported();
    if (this.onImport != null) {
      this.onImport();
    }
  }

  Map<String, dynamic> export([List<String> fields]) {
    var data = {};
    var getters = this.getters;

    getters.forEach((name, getter) {
      data[name] = getter();
    });

    return data;
  }

  async.Future destroy() {
    return this.http.delete(this.url).then((response) {
      this.clearError();
      this.import(response.data);
      return response.data;
    }).catchError((err) {
      return new async.Future.error(this.parseError(err));
    }, test: (e) => e is ng.HttpResponse);
  }

  async.Future send(String method, String url, List<String> fields) {
    var methodFunc;

    var data = this.export(fields);

    if (method == 'post') {
      methodFunc = this.http.post;
    }
    else if (method == 'put') {
      methodFunc = this.http.put;
    }
    else {
      throw new ArgumentError('Unkown method');
    }

    return methodFunc(url, data).then((response) {
      this.clearError();
      this.import(response.data);
      return response.data;
    }).catchError((err) {
      return new async.Future.error(this.parseError(err));
    }, test: (e) => e is ng.HttpResponse);
  }

  async.Future save([List<String> fields]) {
    return this.send('put', this.url, fields);
  }

  async.Future create([List<String> fields]) {
    return this.send('post', this.url, fields);
  }

  void clear() {
    var setters = this.setters;

    setters.forEach((_, setter) {
      setter(null);
    });
  }

  void init() {}
}
