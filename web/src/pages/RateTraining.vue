<template>
    <app-layout>
        <div class="py-5 text-center">
            <p style="font-size: 59px;">✨</p>
            <h2>
                <span>Rate training</span>
            </h2>

            <br>
            <p class="lead">
                It will be cool if you will share with us if you liked your training.
                It will help us to make better trainings in the future.<br>
                Feedback is anonymous.
            </p>
        </div>
        <div class="row justify-content-md-center">
            <div class="col-md-8 order-md-1 l-md">
                <form class="needs-validation" @submit.prevent="rate" novalidate>
                    <div class="row">
                        <div class="col-md-6 mb-3">
                            <label for="rating">Rating</label>
                            <select class="custom-select" size="5" id="rating" v-model="ratingData.rating">
                                <option value="1">⭐</option>
                                <option value="2">⭐ ⭐</option>
                                <option value="3">⭐ ⭐ ⭐</option>
                                <option value="4">⭐ ⭐ ⭐ ⭐</option>
                                <option value="5">⭐ ⭐ ⭐ ⭐ ⭐</option>
                            </select>
                        </div>
                        <div class="col-md-6 mb-3">
                            <div class="form-group">
                                <label for="notes">Comment <small>(visible for trainer)</small></label>
                                <textarea class="form-control" id="notes" rows="6" maxlength="1000"
                                          v-model="ratingData.comment"></textarea>
                            </div>
                        </div>
                    </div>

                    <div class="alert alert-info text-center" role="alert">
                        Training can be rated only once! You can't change your rating later.
                    </div>

                    <button class="btn btn-primary btn-lg btn-block" type="submit">
                        Submit
                        <span v-if="showLoader" class="spinner-grow spinner-grow-sm"
                              style="width: 1.3rem; height: 1.3rem;"
                              role="status" aria-hidden="true"></span>
                    </button>
                </form>
            </div>
        </div>
    </app-layout>
</template>

<script>
import AppLayout from '../layouts/App.vue'
import {submitRating} from "../repositories/ratings";

export default {
    name: "RateTraining",
    params: [],
    components: {
        AppLayout,
    },
    created() {

    },
    data: function () {
        return {
            'ratingData': {
                'rating': 3,
                'comment': '',
            },
            'showLoader': false,
        }
    },
    methods: {
        rate: function () {
            let self = this;
            this.showLoader = true;

            submitRating(
                this.$route.params['trainingID'],
                this.ratingData.rating,
                this.ratingData.comment,
                function () {
                    self.showLoader = false
                    self.$toast.success('Rating submitted!');
                    self.$router.push({name: 'trainingsList'});
                },
                function (err) {
                    self.showLoader = false
                    self.$toast.error("Failed to submit rating");
                    console.error(err)
                }
            )
        }
    }
}
</script>
