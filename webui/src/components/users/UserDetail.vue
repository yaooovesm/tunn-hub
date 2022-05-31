<template>
  <div>
    <el-dialog
        v-model="dialogVisible"
        width="70%"
        custom-class="default-dialog"
        :destroy-on-close="true"
        draggable
    >
      <template #title>
        <div class="title">
          <div class="title-text">用户详情
            <el-tooltip
                effect="dark"
                :content="id"
                placement="bottom-start"
                v-if="id.length>10"
            >
              <el-tag
                  type=""
                  effect="dark"
                  style="transform: translateY(-2px);margin-left: 10px;height: 25px"
              >
                ID&nbsp;{{ id.length > 10 ? id.substring(0, 10) + "..." : id }}
              </el-tag>
            </el-tooltip>
            <el-tag
                type=""
                effect="dark"
                style="transform: translateY(-2px);margin-left: 10px;height: 25px"
                v-else
            >
              ID&nbsp;{{ id }}
            </el-tag>
          </div>

        </div>
      </template>
      <el-scrollbar height="600px">
        <el-row style="padding-top: 30px;padding-bottom: 30px">
          <el-col :span="22" :offset="1">
            <user-card :id="id"/>
          </el-col>
          <el-col :span="22" :offset="1" style="margin-top: 30px">
            <client-card :id="id"/>
          </el-col>
        </el-row>
      </el-scrollbar>
    </el-dialog>
  </div>
</template>

<script>
import ClientCard from "@/components/home/ClientCard";
import UserCard from "@/components/home/UserCard";

export default {
  name: "UserDetail",
  components: {ClientCard, UserCard},

  data() {
    return {
      dialogVisible: false,
      id: ""
    }
  },
  methods: {
    show: function (id) {
      this.dialogVisible = true
      this.id = id
    },
  }
}
</script>

<style scoped>

</style>