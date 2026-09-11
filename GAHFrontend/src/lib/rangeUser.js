import { post } from './request.js'

const rangeUser = async (condition, start, length) => {
    const res = await post('/user/range', {
        begin_table_id: start,
        length: length,
        search_condition: condition
    })
    return res.data
}

export default rangeUser
