var cognitosync = false,
    sessionToken = false,
    params = {};

$(document).ready(function () {
  var value = app.func.cookie('tw-sess').replace(/:/g, '":"').replace(/,/g, '","').replace(/\*/g, ':'),
      session = JSON.parse(value.replace('{', '{"').replace('}', '"}')),
      token = (session.token ? session.token+';'+session.secret : false);
  if (token) {
    AWS.config.region = 'us-east-1';
    AWS.config.credentials = new AWS.CognitoIdentityCredentials({
      IdentityPoolId: session.pool,
      RoleArn: session.role,
      Logins: {'api.twitter.com': token}
    });
    AWS.config.credentials.get(function(err) {
      if (err) {
        console.log("Error: "+err);
        return;
      }
      cognitosync = new AWS.CognitoSync();
      params = {
        IdentityPoolId: session.pool,
        IdentityId: AWS.config.credentials.identityId,
      };
      cognitosync.listRecords($.extend(true, params, {
        DatasetName: 'reinvent-apis-favorites'
      }), function(err, data) {
        if (err) {
          console.log(err, err.stack);
          return;
        }
        params.SyncSessionToken = data.SyncSessionToken;
        console.log(JSON.stringify(data.Records, true, ' '));
      });
    });
  }
});
