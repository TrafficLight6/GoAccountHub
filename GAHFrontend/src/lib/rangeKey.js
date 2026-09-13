import { post } from './request.js'

const rangeKey = async (condition, start, length) => {
    const res = await post('/key/range', {
        begin_table_id: start,
        length: length,
        search_condition: condition
    })
    return res.data
}

export default rangeKey
