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
  <div class="d-flex position-relative">
    <div class="d-flex position-absolute top-0 end-0 mt-3">
      <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
    </div>
  </div>

  <div class="d-flex justify-content-center position-absolute" style="top: 25%; left: 0; width: 100%; height: 100%;">
    <div class="justify-content-between flex-wrap flex-md-nowrap align-items-center mb-2">
      <h1 class="h1">Welcome to WASAPhoto</h1>
      <h1 class="h2 text-center" v-if="this.username">{{ this.username }}</h1>
    </div>
  </div>
  <div class="d-flex justify-content-center position-absolute" style="top: 35%; left: 0; width: 100%; height: 100%;" >
    <div>
      <div class="mb-3 border-bottom"></div>
      <input type="text" id="username" v-model="username" class="form-control"
             placeholder="What's your name?" aria-label="Recipient's username" aria-describedby="basic-addon2">
      <div class="mt-2 text-center">
        <button class="btn" type="button" @click="doLogin">
          <div class="d-flex justify-content-between">
            <svg class="feather mt-1 me-1"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
            <span>Login</span>
          </div>
        </button>
      </div>
    </div>
  </div>


</template>

<style>

</style>