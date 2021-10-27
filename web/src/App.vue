<template>
  <div id="app" v-loading="loading.main">
    <!--筛选&搜索-->
    <div class="flex_row" style="justify-content: space-between;flex-wrap: wrap">
      <el-form :model="formData" :inline="true" size="mini" ref="form1">
        <el-form-item prop="star">
          <el-select v-model="formData.star" clearable filterable @change="do_auto">
            <el-option v-for="(star,i) in stars" :key="'stars'+i" :value="star.Star" :label="star.Star"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="code">
          <el-select v-model="formData.code" clearable filterable @change="do_auto">
            <el-option v-for="(code,i) in codes" :key="'code'+i"
                       :value="code.Code" :label="code.Code"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="ma">
          <el-select v-model="formData.ma" clearable @change="do_auto">
            <el-option value="1" label="有"></el-option>
            <el-option value="2" label="无"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="downed">
          <el-select v-model="formData.downed" clearable @change="do_auto">
            <el-option value="0" label="未下载"></el-option>
            <el-option value="1" label="已下载"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList">搜索</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
      <el-input v-loading="loading.search" placeholder="请输入内容" size="mini"
                style="width: 180px;height: 28px;flex-shrink: 0" v-model="iCode" @keydown.enter.native="do_search">
        <el-button slot="append" icon="el-icon-search" style="height: 28px" @click="do_search"></el-button>
      </el-input>
    </div>

    <el-pagination
            layout="prev, pager, next"
            :background="true"
            :current-page="formData.index"
            :page-size="formData.size"
            :total="formData.count"
            @size-change="size_change"
            @current-change="current_change"
            style="text-align: right;"
    ></el-pagination>

    <el-row :gutter="10">
      <el-col v-for="(p,i) in List" :key="i" class="box" :xs="24" :sm="12" :md="8" :lg="6" :xl="4">
        <i class="el-icon-circle-close del" @click="do_del(p)"></i>
        <el-card class="box-card" shadow="hover" v-loading="loading[p.ID]">
          <div slot="header" class="flex_row box_header">
            <div class="box_title">
              <a :href="busUrl + p.Code" target="_blank" style="color: #333333;">{{ p.Code }}</a>
              (<a :href="busUrl + 'star/' + p.StarCode" target="_blank"
                  :style="{color: p.StarCode?'#333333':'#ff0000'}">{{ p.Star }}</a>)
            </div>
            <el-button type="text" class="rBtn tPadding" @click="do_up(p)">更新</el-button>
          </div>

          <el-image
                  class="box_img" fit="contain"
                  :src="urlFix + '/photos/' + p.Code +'.jpg?' + new Date().getTime()"
                  :preview-src-list="[urlFix + '/photos/' + p.Code +'.jpg?' + new Date().getTime()]"
          ></el-image>

          <div class="flex_row more" style="justify-content: space-between; align-items: center">
            <el-rate v-model="p.Starnum" @change="do_set('st',p)"></el-rate>
            <div>
              <el-button type="text" @click="do_set('ma',p)" title="ma">{{ dic.ma[p.Ima] }}</el-button>
              <el-button type="text" @click="do_set('fa',p)" title="fa">{{ dic.fa[p.Iface] }}</el-button>
              <el-button type="text" @click="do_set('dn',p)" title="dn">{{ dic.downed[p.Downed] }}</el-button>
            </div>

          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-pagination
        layout="prev, pager, next"
        :background="true"
        :current-page="formData.index"
        :page-size="formData.size"
        :total="formData.count"
        @size-change="size_change"
        @current-change="current_change"
        style="text-align: right;"
    ></el-pagination>
  </div>
</template>

<script>

    import axios from 'axios'

    export default {
        name: 'App',
        data() {
            return {
                urlFix: process.env.NODE_ENV === 'production' ? '' : 'http://127.0.0.1:8888',
                busUrl: "",
                loading: {main: false, search: false},
                dic: {ma: ['未知', '无码', '有码'], fa: ['未知', '有颜值', '封面党'], downed: ["未下载", '已下载']},
                iCode: "",
                formData: {
                    downed: null,
                    ma: null,
                    code: null,
                    star: null,
                    size: 48, index: 1, count: 0
                },
                codes: [],
                stars: [],
                List: [
                    // {ID: null, Code: null, Star: null, starcodeStarCode: null, Fname: null, Downed: null, Ima: null, Iface: null, Starnum: null},
                ]
            }
        },
        mounted() {
            this.getList();
        },
        methods: {
            // 表单列表获取
            getList() {
                this.loading.main = true

                axios.request({
                        url: this.urlFix + "/api/getList",
                        params: this.formData
                    }
                ).then(res => {
                    if (res['data']) {
                        res = res['data']
                        this.busUrl = res["busurl"]
                        this.codes = res["codes"]
                        this.stars = res["stars"]
                        this.formData.count = res["count"]
                        this.List = res["fans"]
                    }
                    this.loading.main = false
                }).catch(err => {
                    this.loading.main = false
                    console.log(err)
                })
            },
            do_auto() {
                this.formData.index = 1
                this.getList()
            },

            do_search() {
                if (!this.iCode) {
                    return
                }
                this.loading.search = true
                axios.request({url: this.urlFix + "/api/search", params: {c: this.iCode}}).then(res => {
                    this.loading.search = false
                    let data = res['data']
                    if (data === 'has') {
                        this.$message.warning("番号已存在")
                    } else if (data === 'ok') {
                        this.$message.success("获取成功")
                    }
                    this.iCode = ""
                    this.getList()
                }).catch(err => {
                    this.loading.search = false
                    this.$message.error("获取失败" + err.response.data)
                    console.log(err)
                })
            },

            do_set(st, p) {
                let data = {
                    t: st,
                    id: p.ID,
                    flag: p.Starnum
                }
                switch (st) {
                    case 'st':
                        data.flag = p['Starnum']
                        break
                    case 'ma':
                        data.flag = p['Ima'] === 2 ? 0 : (p['Ima'] + 1)
                        break
                    case 'fa':
                        data.flag = p['Iface'] === 2 ? 0 : (p['Iface'] + 1)
                        break
                    case 'dn':
                        data.flag = p['Downed'] === 1 ? 0 : 1
                        break
                }
                this.$set(this.loading, p.ID, true)
                axios.request({url: this.urlFix + "/api/set", params: data}).then(() => {
                    this.$set(this.loading, p.ID, false)
                    switch (st) {
                        case 'st':
                            p['Starnum'] = data.flag
                            break
                        case 'ma':
                            p['Ima'] = data.flag
                            break
                        case 'fa':
                            p['Iface'] = data.flag
                            break
                        case 'dn':
                            p['Downed'] = data.flag
                            break
                    }
                    this.$message.success("标记成功")
                }).catch(err => {
                    this.$set(this.loading, p.ID, false)
                    this.$message.error("标记失败" + err.response.data)
                    console.log(err)
                })
            },


            do_up(p) {
                this.$set(this.loading, p.ID, true)
                axios.request({url: this.urlFix + "/api/search", params: {c: p.Code, u: 1}}).then(() => {
                    this.$set(this.loading, p.ID, false)
                    this.$message.success("更新成功")
                    this.getList()
                }).catch(err => {
                    this.$set(this.loading, p.ID, false)
                    this.$message.error("更新失败" + err.response.data)
                    console.log(err)
                })
            },

            do_del(p) {
                this.$confirm('此操作将永久删除该番号, 是否继续?', '提示', {type: 'warning'}).then(() => {

                    this.$set(this.loading, p.ID, true)
                    axios.request({url: this.urlFix + "/api/del", params: {qazxsw: p.ID}}
                    ).then(() => {
                        this.$set(this.loading, p.ID, false)
                        this.$message.success("删除成功")
                        this.getList()
                    }).catch(err => {
                        this.$set(this.loading, p.ID, false)
                        this.$message.error("删除失败" + err.response.data)
                        console.log(err)
                    })

                }).catch(() => {
                });
            },

            resetForm() {
                this.$refs['form1'].resetFields();
                this.getList()
            },
            size_change(val) {
                this.formData.size = val
                this.getList()
            },
            current_change(val) {
                this.formData.index = val;
                this.getList();
            },
        }
    }
</script>

<style>
  .box {
    position: relative;
    margin-bottom: 10px;
  }

  .del {
    display: none !important;
    position: absolute;
    right: 50%;
    bottom: -8px;
    color: #00000010;
  }

  .box:hover .del {
    display: inline-block !important;
  }

  .box_header {
    justify-content: space-between;
  }

  .box_title {
    word-break: keep-all;
    white-space: pre;
    overflow: hidden
  }

  .box_img {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 300px;
  }

  .rBtn {
    float: right;
  }

  .red {
    color: #f56c6c !important;
  }

  .tPadding {
    padding: 3px 0 !important;
  }

  .more {
    margin: 0 10px;
  }

  .more .el-rate {
    display: inline-block;
  }

  .flex_row {
    display: flex;
    flex-direction: row;
  }

  .el-card__header {
    padding: 15px 15px !important;
  }

  .el-card__body {
    padding: 0 !important;
  }

  .el-form-item {
    margin-bottom: 10px !important;
  }

  .el-pagination {
    margin-bottom: 10px;
    padding: 0;
  }
  .btn-next{
    margin-right: 0!important;
  }
</style>
