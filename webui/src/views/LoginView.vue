<script>
export default {
  data: function() {
    return {
      error: null,
      username: "",
      ID: 0
    }
  },
  methods: {
    async doLogin() {
      if (this.username === "") {
        this.errormsg = "Username cannot be empty.";
      } else {
        try {
          let response = await this.$axios.put("/session", { username: this.username })
          localStorage.setItem("ID", response.data);
          localStorage.setItem("username", this.username);
          this.$router.push({ path: '/stream' })
        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.err = "Username should has a length between 3 - 16";
          } else if (e.response && e.response.status === 500) {
            this.err = "An internal error occurred, please try again later.";
          } else {
            this.err = e.toString();
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
    <h1 class="h2">Welcome to WASAPhoto</h1>
  </div>
  <div class="input-group mb-3">
    <input type="text" id="username" v-model="username" class="form-control"
           placeholder="Please, insert your username." aria-label="Recipient's username"
           aria-describedby="basic-addon2">
    <div class="input-group-append">
      <button class="btn btn-outline-success" type="button" @click="doLogin">Login</button>
    </div>
  </div>
  <ErrorMsg v-if="err" :msg="err"></ErrorMsg>
</template>

<style>

</style>