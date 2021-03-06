module.exports = function(grunt) {
  // configure the tasks
  grunt.initConfig({
    copy: {
      main: {
        files: [
          {cwd: 'backend/', src: '**', dest: 'build/', expand: true },
          {src: 'assets/*', dest: 'build/static/', expand: true },
          {src: 'resources/*', dest: 'build/', expand: true },
          {src: 'templates/*', dest: 'build/resources/', expand: true },
          {cwd: 'js/src', src: '*.js',
           dest: 'build/static/js/', expand: true },

        ],
      },
    },

    clean: {
      main: {
        src: [ 'build' ],
      },
    },

    sass: {
      main: {
        files: {
          'build/static/css/style.css': 'stylesheets/style.scss'
        }
      }
    },

    autoprefixer: {
      build: {
        expand: true,
        cwd: 'build/static/css/',
        src: [ '*.css' ],
        dest: 'build/static/css/'
      }
    },

    bower: {
      main: {
        dest: 'build/static/deps/',
        options: {
          packageSpecific: {
            bootstrap: {
              dest: 'build/static/deps/fonts',
              css_dest: 'build/static/deps/bootstrap'
            }
          }
        }
      }
    },

    shell: {
      devserver: {
        command: 'dev_appserver.py build/',
        options: {
          async: true,
        },
      },
    },

    watch: {
      stylesheets: {
        files: 'stylesheets/**',
        tasks: [ 'css' ]
      },
      copy: {
        files: [ '{js/src,resources,assets,templates,backend}/**' ],
        tasks: [ 'copy' ]
      }
    },
  });

  // load the tasks
  grunt.loadNpmTasks('grunt-autoprefixer');
  grunt.loadNpmTasks('grunt-bower');
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-copy');
  grunt.loadNpmTasks('grunt-contrib-sass');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-shell-spawn');

  // define the tasks
  grunt.registerTask(
    'css', '',
    [ 'sass', 'autoprefixer' ]
  );
  grunt.registerTask(
    'build',
    'Compiles all of the assets and copies the files to the build directory.',
    [ 'clean', 'bower', 'copy', 'css' ,'bower']
  );
  grunt.registerTask(
    'dev',
    'Builds, runs the dev server, and watches for updates.',
    [ 'build', 'shell:devserver', 'watch']
  );
};
