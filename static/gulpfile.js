/* global require */
/* jshint esversion:6 */
const batch = require('gulp-batch');
const clean = require('gulp-clean');
const concat = require('gulp-concat');
const gulp = require('gulp');
const htmlmin = require('gulp-htmlmin');
const minify = require('gulp-minify');
const rename = require('gulp-rename');
const sourcemaps = require('gulp-sourcemaps');
const strip = require('gulp-strip-comments');
const stylus = require('gulp-stylus');
const watch = require('gulp-watch');

const BUILD_DIRECTORY_BASE = './build';

gulp.task('js', gulp.parallel(function(done) {
    'use strict';

    return gulp.src('./js/ng/**/*.js')
        .pipe(sourcemaps.init())
        .pipe(concat('ng.js'))
        .pipe(sourcemaps.write('.'))
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/js'));
}, function(done) {
    'use strict';

    return gulp.src('./js/*.js')
        .pipe(sourcemaps.init())
        .pipe(sourcemaps.write('.'))
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/js'));
}));

gulp.task('html', gulp.parallel(function(done) {
    'use strict';

    return gulp.src(['./js/ng/pages/**/*.html', './js/ng/components/**/*.html'])
        .pipe(htmlmin({
            collapseWhitespace: true,
            collapseBooleanAttributes: true,
            removeComments: true,
        }))
        .pipe(rename({ dirname: '' }))
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/html/'));
}, function(done) {
    'use strict';

    return gulp.src('./html/*.html')
        .pipe(htmlmin({
            collapseWhitespace: true,
            collapseBooleanAttributes: true,
            removeComments: true,
        }))
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/'));
}));

gulp.task('css', gulp.parallel(function(done) {
    'use strict';

    return gulp.src(['./fonts/*.css'])
        .pipe(concat('fonts.css'))
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/css'));
}, gulp.parallel(function(done) {
    'use strict';

    return gulp.src(['./css/*.styl', './../shared/*.styl'])
        .pipe(concat('main.styl'))
        .pipe(stylus())
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/css'));
})));

gulp.task('img', () => {
    'use strict';

    return gulp.src('./img/**/*')
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/img'));
});

gulp.task('fonts', () => {
    'use strict';

    return gulp.src(['./fonts/*.eot', './fonts/*.svg', './fonts/*.ttf', './fonts/*.woff', './fonts/*.woff2'])
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/fonts'));
});

gulp.task('clean', () => {
    'use strict';

    return gulp.src(BUILD_DIRECTORY_BASE + '/', { read: false })
        .pipe(clean({ force: true }));
});

gulp.task('copy', () => {
    'use strict';

    return gulp.src('./copy/**/*')
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE));
});

gulp.task('watch', gulp.parallel(
    function(done) {
        gulp.watch('./js/**/*.js', gulp.series('js'));
        gulp.watch(['./css/*.styl', './../shared/*.styl'], gulp.series('css'));
        gulp.watch(['./html/*.html', './include/*.html', './js/ng/pages/**/*.html', './js/ng/components/**/*.html'], gulp.series('html'));
        gulp.watch('./img/**/*', gulp.series('img'));
        gulp.watch('./fonts/*', gulp.series('fonts'));
    }
));

gulp.task('release', gulp.parallel(function(done) {
    'use strict';

    return gulp.src([BUILD_DIRECTORY_BASE + '/assets/js/*.js', '!' + BUILD_DIRECTORY_BASE + '/assets/js/*min.js'])
        .pipe(strip())
        .pipe(minify({
            ext: {
                min: '.js'
            },
            noSource: true,
            mangle: false,
            compress: false
        }))
        .pipe(gulp.dest(BUILD_DIRECTORY_BASE + '/assets/js'));
}, function(done) {
    'use strict';

    return gulp.src(BUILD_DIRECTORY_BASE + '/assets/**/*.map', { read: false })
        .pipe(clean({ force: true }));
}));

gulp.task('default', gulp.parallel('css', 'js', 'html', 'img', 'fonts', 'copy'));
gulp.task('start', gulp.series('default', 'watch'));