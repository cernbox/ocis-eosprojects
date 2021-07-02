<template>
  <div class="uk-flex uk-flex-column uk-flex-middle uk-margin-xlarge-top">
    <h1>
      My Projects
    </h1>

    <form v-on:submit.prevent="getProjects()" action="#" method="get">

      <oc-button variation="primary" class="uk-width-1-1 uk-margin-top">
        Get
      </oc-button>
    </form>

    <p class="uk-text-lead" v-if="projects" v-text="projects" />
  </div>
</template>

<script>
export default {
  name: 'App',
  data: function () {
    return {
    }
  },
  computed: {
    projects () {
      return this.$store.getters['EosProjects/projects']
    }
  },
  methods: {
    getProjects () {
      this.$store.dispatch('EosProjects/getProjects')
    }
  }
}
</script>
