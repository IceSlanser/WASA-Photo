<script>
export default {
  data: function() {
    return {
      error: null,
      username: "",
      ID: 0,
    }
  },
  methods: {
    async doLogin() {
      if (this.username === "") {
        this.error = "Username cannot be empty.";
      } else {
        this.error = null
        try {
          let response = await this.$axios.put("/session", { username: this.username })
          await localStorage.setItem("ID", response.data);
          await localStorage.setItem("username", this.username);
          this.$router.push({ path: '/stream' })
        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Username should has a length between 3 - 16";
          } else if (e.response && e.response.status === 500) {
            this.error = "An internal error occurred, please try again later.";
          } else {
            this.error = e.toString();
          }
        }
      }
    }
  }
}
</script>

<template>
  <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 ">
    <h1 class="h2" v-if="this.username">Welcome to WASAPhoto {{this.username}}</h1>
    <h1 class="h2" v-else>Welcome to WASAPhoto </h1>
  </div>
  <div class="input-group mb-2">
    <input type="text" id="username" v-model="username" class="form-control"
           placeholder="Please, insert your username." aria-label="Recipient's username"
           aria-describedby="basic-addon2">
    <div class="input-group-append">
      <button class="btn btn-outline-success" type="button" @click="doLogin">Login</button>
    </div>
  </div>
  <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
</template>

<style>

</style>