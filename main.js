var GitHub = require('github-api');
var menubar = require('menubar');
var mb = menubar()

// basic auth
var gh = new GitHub({
  token: process.env.GITHUB_ACCESS_TOKEN
});

var me = gh.getUser(); // no user specified defaults to the user for whom credentials were provided
me.listRepos(function (err, repos) {
  // do some stuff
  console.log(JSON.stringify(repos));
});

mb.on('ready', function ready() {
  console.log('app is ready')
})



// var clayreimann = gh.getUser('clayreimann');
// clayreimann.listStarredRepos(function(err, repos) {
//    // look at all the starred repos!
// });
