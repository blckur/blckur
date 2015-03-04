library loading;

class Loading {
  bool loading;

  bool setLoading() {
    if (this.loading == true) {
      return false;
    }
    this.loading = true;
    return true;
  }

  void clearLoading() {
    this.loading = false;
  }
}
