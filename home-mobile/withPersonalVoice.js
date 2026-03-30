const { withInfoPlist } = require('@expo/config-plugins');

module.exports = function withPersonalVoice(config) {
  return withInfoPlist(config, (config) => {
    // 1. Permission to show the "Allow this app to use your voice" popup
    config.modResults.NSSpeechRecognitionUsageDescription = 
      "The AI uses your Personal Voice to give you a natural experience.";
    
    // 2. Permission for the microphone (needed for AI interaction)
    config.modResults.NSMicrophoneUsageDescription = 
      "Talk to your AI assistant.";

    return config;
  });
};