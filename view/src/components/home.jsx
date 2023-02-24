import React, {useState} from "react"
import "./home.css"
import {Button, Card, Col, Divider, Input, List, message, Row, Select} from "antd"
import axios from "axios"
import {useNavigate} from "react-router"

export default function () {
    let nav = useNavigate()
    let [surname, setSurname] = useState("李")
    let [surnameInfo, setSurnameInfo] = useState({})
    let [data, setData] = useState({})
    let [page, setPage] = useState(1)
    let [length, setLength] = useState(2)
    let [rhyme, setRhyme] = useState("rhyme")
    let [tone, setTone] = useState(0)
    let [pagesize, setPagesize] = useState(10)
    let [loading, setLoading] = useState(false)

    const loadSurnameInfo = () => {
        axios.get(`/word/pinyin/${surname}`).then(data => {
            setSurnameInfo(data.data.body)
        })
    }

    useState(() => {
        loadSurnameInfo()
    })

    const loadWords = (page, pagesize) => {
        setLoading(true)
        loadSurnameInfo()
        let url = `/words/${rhyme}/${surname}?pagesize=${pagesize}&page=${page}&length=${length}&rhyme=${rhyme}&tone=${tone}`
        axios.get(url).then(({data}) => {
            setLoading(false)
            if (data.status_code !== 200) {
                return message.error(`发生错误: ${data.message}`)
            }

            setData(data.body)
        })
    }

    const wordHighlight = (contents, word) => {

        return contents.replaceAll(word, `<span class="keyword">${word}</span>`)
        let t = contents
        let c = []
        let p = 0
        while (p < contents.length) {
            let index = t.indexOf(word, p)
            if (index === -1) {
                continue
            }

            c.push({
                text: t.substring(p, index),
                keyword: word,
            })
            p = index + word.length
        }

        return c
    }

    return <>
        <div className="inputer">
            <Row align="middle">
                <Col span={6} style={{textAlign: "center"}}>
                    <h1>Name From Poetries</h1>
                </Col>
                <Col span={12} style={{textAlign: "center"}}>
                    <Input.Group compact>
                        <Input
                            value={surname}
                            style={{width: "calc(100% - 350px)", textAlign: "left"}}
                            placeholder="输入姓氏"
                            onChange={e => setSurname(e.target.value)}
                        />
                        <Select placeholder="长度" defaultValue="2" value={length} onChange={v => setLength(v)}>
                            <Option value={1}>单字</Option>
                            <Option value={2}>双字</Option>
                            <Option value={3}>三字</Option>
                            <Option value={4}>四字</Option>
                            <Option value={5}>不限</Option>
                        </Select>
                        <Select placeholder="声调" value={tone} onChange={v => setTone(v)}>
                            <Option value={0}>不限</Option>
                            <Option value={1}>阴平</Option>
                            <Option value={2}>阳平</Option>
                            <Option value={3}>上声</Option>
                            <Option value={4}>去声</Option>
                        </Select>
                        <Select placeholder="押韵" defaultValue="rhyme" value={rhyme} onChange={v => setRhyme(v)}>
                            <Option value="rhyme">压韵母</Option>
                            <Option value="rhyme_tone">押韵母+声调</Option>
                            <Option value="random">不押韵</Option>
                        </Select>
                        <Button type="primary" onClick={() => {
                            setPage(1)
                            loadWords(1, 10)
                        }}>查询</Button>
                    </Input.Group>
                </Col>
            </Row>
        </div>
        <div style={{padding: "20px"}}>
            {
                page > 1 ? "" : <Row>
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
            }
            <Row>
                <Col span={24}>
                    <List
                        loading={loading}
                        itemLayout="vertical"
                        size="large"
                        pagination={{
                            onChange(page, pagesize) {
                                setPage(() => {
                                    loadWords(page, pagesize)
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
                        dataSource={data.words}
                        renderItem={(item) => <List.Item
                            key={<span>{item.label}</span>}
                            actions={[
                                <span>在古诗词中的出现次数：<a
                                    onClick={() => nav(`/detail?word=${item.word}`)}>{item.use_count}</a></span>,
                            ]}
                        >
                            <List.Item.Meta
                                title={<span>{surname} {item.word}({item.pinyin_tone})</span>}
                                description={<span>出处：{item.label}</span>}
                            />
                            <span dangerouslySetInnerHTML={{
                                __html: item.contents.replaceAll(
                                    item.word, `<span class="keyword">${item.word}</span>`,
                                ),
                            }}></span>
                        </List.Item>}
                    />
                </Col>
            </Row>
        </div>
    </>
}
