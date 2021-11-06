<template>
  <app-layout>
    <div class="py-5 text-center">
      <h2>Your ratings ✨</h2>
      <p class="lead">Check if your customers are happy with your trainings.</p>
    </div>
    <br><br>
    <table class="table">
      <thead>
      <tr>
        <th scope="col">Date</th>
        <th scope="col">Rating</th>
        <th scope="col">Comment</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(rating, idx) in ratings" :key="idx" v-bind:class="{ 'text-danger': rating.rating === 1  }">
        <th scope="row">{{ formatDate(rating.date) }}</th>
        <td><span v-for="n in rating.rating" :key="n">⭐</span></td>
        <td><span v-if="rating.comment">{{ rating.comment }}</span><span class="text-secondary" v-if="!rating.comment">(no comment left)</span></td>
      </tr>
      </tbody>
    </table>
  </app-layout>
</template>

<script>
import AppLayout from '../layouts/App.vue'
import {getRatings} from "../repositories/ratings";
import {formatDate} from "../date";

export default {
  components: {
    AppLayout,
  },
  data: function () {
    return {
      'ratings': null,
    }
  },
  mounted() {
    let self = this
    getRatings(function (ratings) {
      self.ratings = ratings
    })
  },
  methods: {
      formatDate,
  }
}
</script>

