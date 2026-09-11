import { post } from './request.js'

const rangeAdmin = async (condition, start, length) => {
    const res = await post('/admin/range', {
        begin_table_id: start,
        length: length,
        search_condition: condition
    })
    return res.data
}

export default rangeAdmin
