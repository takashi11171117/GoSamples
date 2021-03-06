window.addEventListener('load', function () {
  document.getElementById('sign-out').onclick = function () {
    firebase.auth().signOut();
  };

  // [START UIconfig_variable]
  // FirebaseUI config.
  var uiConfig = {
    signInSuccessUrl: '/',
    signInOptions: [
      firebase.auth.GoogleAuthProvider.PROVIDER_ID,
      // firebase.auth.FacebookAuthProvider.PROVIDER_ID,
      // firebase.auth.TwitterAuthProvider.PROVIDER_ID,
      // firebase.auth.GithubAuthProvider.PROVIDER_ID,
      // firebase.auth.EmailAuthProvider.PROVIDER_ID,
      // firebase.auth.PhoneAuthProvider.PROVIDER_ID
    ],
    tosUrl: ''
  };
  // [END UIconfig_variable]

  // [START auth_request]
  firebase.auth().onAuthStateChanged(function (user) {
    if (user) {
      // User is signed in.
      document.getElementById('sign-out').hidden = false;
      document.getElementById('post-form').hidden = false;
      document.getElementById('account-details').textContent =
        'Signed in as ' + user.displayName + ' (' + user.email + ')';
      user.getIdToken().then(function (accessToken) {
        document.getElementById('token').value = accessToken;
      });
    } else {
      var ui = new firebaseui.auth.AuthUI(firebase.auth());
      ui.start('#firebaseui-auth-container', uiConfig);
      document.getElementById('sign-out').hidden = true;
      document.getElementById('post-form').hidden = true;
      document.getElementById('account-details').textContent = '';
    }
  }, function (error) {
    console.log(error);
    alert('Unable to log in: ' + error)
  });
  // [END auth_request]
});