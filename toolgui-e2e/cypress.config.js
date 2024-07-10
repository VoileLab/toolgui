const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    projectId: "rtf9hd",
    baseUrl: 'http://localhost:3000',
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
});
