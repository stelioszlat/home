import { AppleSpeech } from '@react-native-ai/apple';

// const requestVoiceAccess = async () => {
//   const status = await AppleSpeech.requestPersonalVoiceAuthorization();
//   if (status === 'authorized') {
//     console.log("Success! We can use the clone.");
//   }
// };

const getDefaultVoice = async () => {
  const voices = await AppleSpeech.getVoices();
  
  // Look for the voice marked as 'personal'
  const personalVoice = voices.find(v => v.name.includes("Zlat try hard English"));
  
  return personalVoice?.identifier;
};

const speakAIResponse = async (text: string) => {
  const voiceId = await getDefaultVoice();
  
  await AppleSpeech.generate(text, {
    voice: voiceId,
  });
};

module.exports = {
    getDefaultVoice: getDefaultVoice,
    speakAIResponse: speakAIResponse
}