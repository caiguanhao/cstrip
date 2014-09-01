module.exports = function(grunt) {

  grunt.initConfig({
    clean: {
      dist: 'dist',
      tpl: 'dist/templates.js'
    },
    copy: {
      data: {
        files: [{
          expand: true,
          cwd: 'data/',
          dest: 'dist/',
          src: '*'
        }]
      },
      html: {
        files: {
          'dist/index.html': 'public/index.production.html'
        }
      },
      fonts: {
        files: [{
          expand: true,
          cwd: 'public/vendor/',
          dest: 'dist/',
          src: 'fonts/*'
        }]
      },
      images: {
        files: [{
          expand: true,
          cwd: 'public/',
          dest: 'dist/',
          src: 'images/*'
        }]
      }
    },
    cssmin: {
      css: {
        files: {
          'dist/assets/app.css': [
            'public/vendor/css/bootstrap.css',
            'public/vendor/css/bootstrap-theme.css',
            'public/css/main.css'
          ]
        }
      }
    },
    rename: {
      assets: {
        options: {
          callback: function(befores, afters) {
            var publicdir = require('fs').realpathSync('dist');
            var path = require('path');
            var index = grunt.file.read('dist/index.html'), before, after;
            for (var i = 0; i < befores.length; i++) {
              before = path.relative(publicdir, befores[i]);
              after = path.relative(publicdir, afters[i]);
              index = index.replace(before, after);
            }
            grunt.file.write('dist/index.html', index);
          }
        },
        files: [
          {
            src: [
              'dist/assets/*.css',
              'dist/assets/*.js'
            ]
          }
        ]
      }
    },
    uglify: {
      js: {
        files: {
          'dist/assets/app.js': [
            'public/vendor/js/angular.js',
            'public/vendor/js/angular-route.js',
            'public/js/main.js',
            'public/js/production.js',
            'dist/templates.js'
          ]
        }
      }
    },
    yaat: {
      CS: {
        options: {
          keyNameCallback: function(name) {
            return name.replace(/^public/, '');
          }
        },
        files: {
          'dist/templates.js': 'public/views/*.html'
        }
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-copy');
  grunt.loadNpmTasks('grunt-contrib-cssmin');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-rename-assets');
  grunt.loadNpmTasks('grunt-yet-another-angular-templates');

  grunt.registerTask('default', [
    'clean:dist',
    'copy',
    'cssmin',
    'yaat',
    'uglify',
    'clean:tpl',
    'rename',
    'compress'
  ]);

  grunt.registerTask('compress', 'Compress assets files', function() {
    var finish = this.async();
    var fs = require('fs');
    var exec = require('child_process').exec;
    exec('gzip -f1k assets/* *.json', {
      cwd: fs.realpathSync('dist')
    }, function(error, stdout, stderr) {
      if (stderr) grunt.fail.fatal(stderr);
      if (error) grunt.fail.fatal(error);
      grunt.log.ok('Asset files compressed.')
      finish();
    });
  });

};
