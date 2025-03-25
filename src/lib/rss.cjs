const Parser = require('rss-parser');
const parser = new Parser();

module.exports = {
  rssGet: async (url) => {
    let feed = await parser.parseURL(url);
    return feed.items
  }
}