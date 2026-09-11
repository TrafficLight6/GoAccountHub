<template>
    <el-card>
        <template #header>
            操作
        </template>
        <p><el-input v-model="searchForm.username" placeholder="用户名"></el-input></p>
        <p><el-input v-model="searchForm.uu_hash" placeholder="用户UUHash"></el-input></p>
        <template #footer>
            <div style="text-align: right">
                <el-button type="primary" @click="handleSearch">查询</el-button>
                <el-button type="warning" @click="handleReset">重置</el-button>
                <el-button type="success" @click="handleAdd">添加用户</el-button>
            </div>
        </template>
    </el-card>
    <br>
    <template v-if="selectedIds.length > 0">
        <el-card>
            <template #header>
                批量操作
            </template>
            <el-text type="primary" size="large" style="text-align: center">已选中 {{ selectedIds.length }} 条</el-text>
            <p>
                <el-button type="primary" @click="handleCheckboxCancel">取消选中</el-button>
                <el-button type="danger" @click="handleBatchDelete">删除选中用户</el-button>
            </p>
        </el-card>
        <br>
    </template>
    <el-card>
        <template #header>
            用户列表（共 {{ userCount }} 位）
        </template>
        <div style="height: 600px">
            <el-auto-resizer>
                <template #default="{ height, width }">
                    <el-table-v2 :columns="columns" :data="userList" :width="width" :height="height"
                        row-key="ID" :footer-height="noMore ? 32 : 0" fixed @end-reached="handleEndReached">
                        <template #footer>
                            <el-text type="success" v-if="noMore" size="large">已加载全部 {{ userList.length }} 条</el-text>
                        </template>
                    </el-table-v2>
                </template>
            </el-auto-resizer>
        </div>
    </el-card>

    <AddUserDialogFrom v-model="dialogVisible" @submit="handleAddSubmit" />
    <EditUserDialog v-model="editVisible" :user="editRow" @submit="handleEditSubmit" />
</template>
<script setup>
import { post, put, del } from '.././../../lib/request.js'
import rangeUser from '.././../../lib/rangeUser.js'
import AddUserDialogFrom from '../../customize/AddUserDialogFrom.vue'
import EditUserDialog from '../../customize/EditUserDialog.vue'
import { h, ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElCheckbox, ElButton, ElIcon, ElMessage, ElMessageBox } from 'element-plus'
import { DocumentCopy } from '@element-plus/icons-vue'

const router = useRouter()
const adminInfo = ref({})
const userList = ref([])
// 用户总数（/user/count）
const userCount = ref(0)

const PAGE_SIZE = 100
const loading = ref(false)
// 后端是否已无更多数据
const noMore = ref(false)

// 筛选条件表单（提交后由后端筛选）
const searchForm = reactive({
    username: '',
    uu_hash: '',
})

// 当前生效的筛选条件，随每次 range 请求发给后端
const searchCondition = ref({})

// 添加用户弹窗
const dialogVisible = ref(false)

// 编辑用户弹窗
const editVisible = ref(false)
const editRow = ref(null)

const canGetUser = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_operate_user)
})

// 勾选状态：以行 ID 记录
const selectedIds = ref([])
const isSelected = (id) => selectedIds.value.includes(id)
// 全选判定基于当前数据
const allSelected = computed(
    () => userList.value.length > 0 && userList.value.every((row) => isSelected(row.ID))
)

const toggleRow = (id) => {
    selectedIds.value = isSelected(id)
        ? selectedIds.value.filter((item) => item !== id)
        : [...selectedIds.value, id]
}

const toggleAll = () => {
    const shownIds = userList.value.map((item) => item.ID)
    selectedIds.value = allSelected.value
        ? selectedIds.value.filter((id) => !shownIds.includes(id))
        : [...new Set([...selectedIds.value, ...shownIds])]
}

// 取消全部选中
const handleCheckboxCancel = () => {
    selectedIds.value = []
}

// 查询：把筛选条件交给后端，重新从第一页取
const handleSearch = () => {
    searchCondition.value = {
        username: searchForm.username.trim(),
        uu_hash: searchForm.uu_hash.trim(),
    }
    range(1)
}

// 重置：清空筛选条件并重新从第一页取全部
const handleReset = () => {
    searchForm.username = ''
    searchForm.uu_hash = ''
    searchCondition.value = {}
    range(1)
}

// 打开添加用户弹窗
const handleAdd = () => {
    dialogVisible.value = true
}

// 提交添加：调用 /user/add，成功后刷新列表
const handleAddSubmit = async (form) => {
    try {
        await post('/user/add', {
            username: form.username.trim(),
            password: form.password,
            meta_data: form.meta_data,
        }, { autoRedirect401: false })
        ElMessage.success('添加用户成功')
        dialogVisible.value = false
        refresh()
        fetchCount()
    } catch {
        // 失败提示已由 request 封装统一弹出
    }
}

// 执行删除并同步列表（后端每次只接收一个 uu_hash，故逐条调用）
const deleteUsers = async (rows) => {
    const results = await Promise.allSettled(
        rows.map((item) => del('/user/delete', { uu_hash: item.UUHash }, { autoRedirect401: false, showError: false }))
    )
    const successIds = rows.filter((_, index) => results[index].status === 'fulfilled').map((item) => item.ID)
    const failed = results.filter((result) => result.status === 'rejected')
    // 成功的行取消勾选
    selectedIds.value = selectedIds.value.filter((id) => !successIds.includes(id))
    if (failed.length === 0) {
        ElMessage.success(`删除成功，共 ${successIds.length} 条`)
    } else if (successIds.length === 0) {
        ElMessage.error(failed[0].reason?.message || '删除失败')
    } else {
        ElMessage.warning(`成功 ${successIds.length} 条，失败 ${failed.length} 条：${failed[0].reason?.message ?? '未知原因'}`)
    }
    if (successIds.length > 0) {
        refresh()
        fetchCount()
    }
}

// 单个删除：确认后删除该行
const handleDelete = async (row) => {
    try {
        await ElMessageBox.confirm(`确定删除用户"${row.Username}"吗？`, '删除确认', {
            type: 'warning',
            confirmButtonText: '确定',
            cancelButtonText: '取消',
        })
    } catch {
        // 取消删除
        return
    }
    deleteUsers([row])
}

// 批量删除：确认后删除所有选中行
const handleBatchDelete = async () => {
    const rows = userList.value.filter((item) => selectedIds.value.includes(item.ID))
    if (rows.length === 0) return
    try {
        await ElMessageBox.confirm(`确定删除选中的 ${rows.length} 条用户吗？`, '删除确认', {
            type: 'warning',
            confirmButtonText: '确定',
            cancelButtonText: '取消',
        })
    } catch {
        // 取消删除
        return
    }
    deleteUsers(rows)
}

// 打开编辑用户弹窗
const handleEdit = (row) => {
    editRow.value = row
    editVisible.value = true
}

// 提交编辑：调用 /user/edit，成功后刷新列表
const handleEditSubmit = async (form) => {
    try {
        await put('/user/edit', {
            uu_hash: form.uu_hash,
            username: form.username.trim(),
            password: form.password,
            meta_data: form.meta_data,
        }, { autoRedirect401: false })
        ElMessage.success('修改用户成功')
        editVisible.value = false
        refresh()
    } catch {
        // 失败提示已由 request 封装统一弹出
    }
}

// 复制文本到剪贴板
const handleCopy = async (text) => {
    try {
        await navigator.clipboard.writeText(text)
        ElMessage.success('已复制')
    } catch {
        ElMessage.error('复制失败')
    }
}

// MetaData 可能很长，单元格内截断显示，hover 看全文
const metaDataCell = ({ rowData }) => h('span', {
    style: 'display: block; overflow: hidden; white-space: nowrap; text-overflow: ellipsis;',
    title: rowData.MetaData ?? '',
}, rowData.MetaData ?? '')

const columns = [
    {
        key: 'selection',
        width: 50,
        headerRenderer: () => h(ElCheckbox, {
            modelValue: allSelected.value,
            'onUpdate:modelValue': toggleAll,
            size: 'large',
        }),
        cellRenderer: ({ rowData }) => h(ElCheckbox, {
            modelValue: isSelected(rowData.ID),
            'onUpdate:modelValue': () => toggleRow(rowData.ID),
            size: 'large',
        }),
    },
    { key: 'ID', dataKey: 'ID', title: 'ID', width: 70 },
    { key: 'Username', dataKey: 'Username', title: '用户名', width: 140 },
    {
        key: 'UUHash',
        title: 'UUHash',
        width: 600,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; align-items: center; gap: 8px;' }, [
            h('span', rowData.UUHash),
            h(ElButton, {
                size: 'small',
                title: '复制 UUHash',
                onClick: () => handleCopy(rowData.UUHash),
            }, () => h(ElIcon, null, () => h(DocumentCopy))),
        ]),
    },
    { key: 'MetaData', title: 'MetaData', width: 300, cellRenderer: metaDataCell },
    {
        key: 'actions',
        title: '操作',
        width: 200,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; gap: 8px;' }, [
            h(ElButton, { type: 'danger', size: 'small', onClick: () => handleDelete(rowData) }, () => '删除用户'),
            h(ElButton, { type: 'primary', size: 'small', onClick: () => handleEdit(rowData) }, () => '编辑用户'),
        ]),
    },
]

const range = async (start = 1) => {
    if (loading.value) return
    loading.value = true
    try {
        const list = await rangeUser(searchCondition.value, start, PAGE_SIZE)
        // refresh() 之后继续下拉时，请求的页可能与已加载数据重叠，按 ID 去重
        const map = new Map()
        for (const row of start === 1 ? list : [...userList.value, ...list]) {
            map.set(row.ID, row)
        }
        userList.value = [...map.values()].sort((a, b) => a.ID - b.ID)
        noMore.value = list.length < PAGE_SIZE
    } catch {
        // 401 已由 request 封装自动跳转登录页
    } finally {
        loading.value = false
    }
}

// 增删改后刷新：按当前已加载的行数重新取。
// 不能只取第一页，否则已加载的那部分行会被丢掉，被操作的行看上去“消失”了
const refresh = async () => {
    if (loading.value) return
    loading.value = true
    try {
        const length = Math.max(PAGE_SIZE, userList.value.length)
        const list = await rangeUser(searchCondition.value, 1, length)
        userList.value = list.sort((a, b) => a.ID - b.ID)
        noMore.value = list.length < length
    } catch {
        // 401 已由 request 封装自动跳转登录页
    } finally {
        loading.value = false
    }
}

// 滚动到底部：后端还有数据时继续请求下一页 100 行
const handleEndReached = () => {
    if (loading.value || noMore.value) return
    range(Math.floor(userList.value.length / PAGE_SIZE) + 1)
}

// 用户总数
const fetchCount = async () => {
    try {
        const res = await post('/user/count')
        userCount.value = res.data.user_count
    } catch {
        // 401 已由 request 封装自动跳转登录页
    }
}

onMounted(async () => {
    try {
        const res = await post('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 已由 request 封装自动跳转登录页，网络错误时不再继续鉴权
        return
    }
    if (!canGetUser.value) {
        router.push('/')
        return
    }
    range()
    fetchCount()
})
</script>
