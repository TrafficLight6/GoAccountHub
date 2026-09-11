import { post } from './request.js'

const rangeCharacter = async (condition, start, length) => {
    const res = await post('/character/range', {
        begin_table_id: start,
        length: length,
        search_condition: condition
    })
    return res.data
}

export default rangeCharacter
