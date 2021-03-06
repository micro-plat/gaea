package pages

const MenuTpl = `
<template>
  <div id="app">
    <nav-menu
      :menus="menus"
      :pwd="pwd"
      ref="NewTap"
    >
      <keep-alive>
      <router-view v-if="$route.meta.keepAlive" @addTab="addTab" @close="close" @setTab="setTab" ></router-view>
      </keep-alive>
      <router-view v-if="!$route.meta.keepAlive" @addTab="addTab" @close="close" @setTab="setTab"></router-view>
       <!--<router-view  @addTab="addTab" @close="close"></router-view>-->
    </nav-menu>
    <!-- Add Form -->
    <el-dialog title="修改密码" width="30%" :visible.sync="dialogAddVisible">
      <el-form :model="updateInfo" :rules="rules" ref="addForm">

        <el-form-item label="请输入原密码" prop="password_old">
          <el-input type="password" v-model="updateInfo.password_old"  ></el-input>
        </el-form-item>

        <el-form-item label="请输入新密码" prop="password">
          <el-input type="password" v-model="updateInfo.password"  ></el-input>
        </el-form-item>

        <el-form-item label="请确认密码" prop="checkPass">
          <el-input type="password" v-model="updateInfo.checkPass"  ></el-input>
        </el-form-item>

      </el-form>
      <div slot="footer" >
        <button class="btn btn-sm btn-primary" @click="resetForm('addForm')">取 消</button>
        <button class="btn btn-sm btn-danger"  @click="add('addForm')">确 定</button>
      </div>

    </el-dialog>
    <!--Add Form -->
  </div>
</template>
 
<script>
  import navMenu from 'nav-menu'; // 引入
  export default {
    name: 'app',
    data () {

      var validatePass = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请输入密码'));
        } else {
          if (this.updateInfo.checkPass !== '') {
            this.$refs.addForm.validateField('checkPass');
          }
          callback();
        }
      };
      var validatePass2 = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请再次输入密码'));
        } else if (value !== this.updateInfo.password) {
          callback(new Error('两次输入密码不一致!'));
        } else {
          callback();
        }
      };
      return {
        menus: [{}],  //菜单数据
        indexUrl: "/",
        dialogAddVisible:false,     //添加表单显示隐藏
        updateInfo:{
          password_old: "",
          password: "",
          checkPass: "",
        },
        rules: {                    //数据验证规则
          password_old: [
            { required: true, message: "请输入原密码", trigger: "blur" }
          ],
          password: [
            { required: true, message: "请输入新密码", trigger: "blur" },
            { validator: validatePass, trigger: 'change' }
          ],
          checkPass: [
            { required: true, message: "请确认密码", trigger: "blur" },
            { validator: validatePass2, trigger: 'change' }
          ],
        },

      }
    },
    components:{ //注册插件
      navMenu
    },
    created(){
      this.getMenu();
    },
    mounted(){
      this.$refs.NewTap.add("首页", this.indexUrl ,{});   //设置默认页面
      document.title = "{{.projectName}} 系统";
    },
    methods:{
      pwd(val){

        this.dialogAddVisible = val;
      },
      resetForm(formName) {
        this.dialogAddVisible = false;
        this.$refs[formName].resetFields();
      },
      add(formName){
        // console.log(this.addData)
        this.$refs[formName].validate((valid) => {
          if (valid) {
            this.$get("/member/changepwd",{
              password_old : this.updateInfo.password_old,
              password : this.updateInfo.password,
            }).then(res=>{
              this.$notify({
                title:'成功',
                message:'修改操作完成',
                type:'success'
              });
              this.dialogAddVisible = false;
              this.$refs[formName].resetFields();
            }).catch(errro=>{
              this.$notify({
                title:'失败',
                message:"原密码错误或密码修改次数超过限制",
                type:'error'
              });
              this.$refs[formName].resetFields();
            })
          } else {
            console.log('error submit!!');
            return false;
          }
        });
      },

      getMenu(){
        this.$get("/member/menuget")
          .then(res => {
            this.menus = res;
          })
          .catch(err => {
           console.log(err)
          });
      },
      //@name 标签名称
      //@path 路由
      //@obj  路由参数 类型：Object
      addTab(name,path,obj){
        this.$refs.NewTap.add(name,path,obj);   //调用组件方法，添加一个页面
      },
      close(v){
         this.$refs.NewTap.closeTab(v);   
      },
      setTab(name,path,obj){
        console.log("outer",name,path,obj);
        this.$refs.NewTap.set(name,path,obj);
      }   
    }
  }
</script>
 
<style scoped>
 
</style>`
