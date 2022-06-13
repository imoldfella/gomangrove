module.exports = {
	globDirectory: '.school',
	globPatterns: [
		'**/*.{jpg,html,jpeg,md,json,css}'
	],
	swDest: '.school/sw.js',
	ignoreURLParametersMatching: [
		/^utm_/,
		/^fbclid$/
	]
};