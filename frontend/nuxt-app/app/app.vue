<template>
  <div class="container mx-auto p-4">
    <h1 class="text-3xl font-bold mb-6">Activity Feed & RAG Chat</h1>
    
    <div class="mb-8 p-4 border rounded-md shadow-sm bg-gray-50">
      <h2 class="text-xl font-semibold mb-2">New Activity</h2>
      <form @submit.prevent="createActivity" class="flex flex-col space-y-2">
        <input
          v-model="title"
          placeholder="Activity Title"
          class="p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <textarea
          v-model="description"
          placeholder="What's on your mind?"
          class="p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows="3"
        ></textarea>
        <button
          type="submit"
          class="bg-blue-600 text-white p-2 rounded-md hover:bg-blue-700 transition-colors"
        >
          Add Activity
        </button>
      </form>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <div>
        <h2 class="text-xl font-semibold mb-2">Activity Feed</h2>
        <div v-if="activities.length === 0" class="text-gray-500 italic">
          No activities yet. Add one!
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="activity in activities"
            :key="activity.ID"
            class="p-4 border rounded-md bg-white shadow-sm"
          >
            <p>{{ activity.title }}</p>
            <p>{{ activity.description }}</p>
            <span class="text-sm text-gray-500 mt-2 block">{{ new Date(activity.CreatedAt).toLocaleString() }}</span>
          </div>
        </div>
      </div>

      <div>
        <h2 class="text-xl font-semibold mb-2">Ask a question (RAG)</h2>
        <div class="border rounded-md p-4 flex flex-col h-96 bg-white shadow-sm">
          <div class="flex-grow overflow-y-auto mb-4 space-y-2">
            <div v-for="(message, index) in chatHistory" :key="index" :class="{'text-right': message.from === 'user'}">
              <span
                :class="{
                  'bg-gray-200 text-black': message.from === 'llm',
                  'bg-blue-500 text-white': message.from === 'user'
                }"
                class="inline-block p-2 rounded-lg"
              >
                {{ message.text }}
              </span>
            </div>
          </div>
          <form @submit.prevent="askQuestion" class="flex space-x-2">
            <input
              v-model="chatQuestion"
              type="text"
              placeholder="Pergunte sobre as atividades..."
              class="flex-grow p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              type="submit"
              class="bg-blue-600 text-white p-2 rounded-md hover:bg-blue-700 transition-colors"
            >
              Send
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const title = ref('');
const description = ref('');
const activities = ref([]);
const chatQuestion = ref('');
const chatHistory = ref([]);

const API_URL = 'http://localhost:8080/api';

const fetchActivities = async () => {
  const response = await fetch(`${API_URL}/activities`);
  activities.value = await response.json();
};

async function createActivity() {
  try {
    const activityData = {
      title: title.value,
      description: description.value,
      date: new Date().toISOString(),
    };
    console.log('Creating activity:', activityData);
    const response = await fetch('http://localhost:8080/api/activities', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json', // <--- Isso é crucial!
      },
      body: JSON.stringify(activityData), // <--- Isso garante que o corpo é JSON
    });
    const data = await response.json();
    console.log('Activity created:', data);
  } catch (error) {
    console.error('Failed to create activity:', error);
  }
}

const askQuestion = async () => {
  if (!chatQuestion.value) return;

  chatHistory.value.push({ text: chatQuestion.value, from: 'user' });

  const response = await fetch(`${API_URL}/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ question: chatQuestion.value }),
  });
  
  const data = await response.json();
  chatHistory.value.push({ text: data.answer, from: 'llm' });
  chatQuestion.value = '';
};

onMounted(() => {
  fetchActivities();
});
</script>

<style>
@import url('https://cdnjs.cloudflare.com/ajax/libs/tailwindcss/2.2.19/tailwind.min.css');
</style>