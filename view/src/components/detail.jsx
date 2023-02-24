import React, {useState} from "react"
import {Card, Col, Divider, List, message, Row} from "antd"
import axios from "axios"
import {useSearchParams} from "react-router-dom"

export default function () {
    let [surnameInfo, setSurnameInfo] = useState({})
    let [data, setData] = useState({})
    let [page, setPage] = useState(1)
    let [pagesize, setPagesize] = useState(10)
    let [loading, setLoading] = useState(false)
    let [params] = useSearchParams()
    let word = params.get("word")

    const loadSurnameInfo = () => {
        axios.get(`/word/pinyin/${word}`).then(data => {
            setSurnameInfo(data.data.body)
        })
    }

    const loadSources = () => {
        setLoading(true)
        let url = `/word/sources/${word}?pagesize=${pagesize}&page=${page}`
        axios.get(url).then(({data}) => {
            setLoading(false)
            if (data.status_code !== 200) {
                return message.error(`发生错误: ${data.message}`)
            }

            setData(data.body)
        })
    }

    useState(() => {
        loadSurnameInfo()
        loadSources()
    })

    return <>
        <div style={{padding: "20px"}}>
            <Row>
                <Col span={24}>
                    <Card type="inner" title={surnameInfo.word}>
                        <p>读音：{surnameInfo.pinyin !== "" ? surnameInfo.pinyin : "-"}</p>
                        <p>声母：{surnameInfo.rhyme_initials !== "" ? surnameInfo.rhyme_initials : "-"}</p>
                        <p>押韵韵母：{surnameInfo.rhyme !== "" ? surnameInfo.rhyme : "-"}</p>
                        <p>押韵音调：{[0, "阴平", "阳平", "上声", "去声"][surnameInfo.rhyme_tone_name]}</p>
                    </Card>
                </Col>
                <Divider/>
            </Row>
            <Row>
                <Col span={24}>
                    <List
                        loading={loading}
                        itemLayout="vertical"
                        size="large"
                        pagination={{
                            onChange(page, pagesize) {
                                setPage(() => {
                                    loadSources(page, pagesize)
                                    setPagesize(pagesize)
                                    return page
                                })
                            },
                            showSizeChanger: true,
                            current: page,
                            total: data.sum_count,
                            pageSize: pagesize,
                            hideOnSinglePage: true,
                            showTotal: total => `共${total}条`,
                        }}
                        dataSource={data.list}
                        renderItem={(item) => <List.Item key={<span>{item.label}</span>}>
                            <List.Item.Meta
                                title={<span>{item.label}</span>}
                            />
                            <span dangerouslySetInnerHTML={{
                                __html: item.contents.replaceAll(
                                    word, `<span class="keyword">${word}</span>`,
                                ),
                            }}></span>
                        </List.Item>}
                    />
                </Col>
            </Row>
        </div>
    </>
}
