"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { useAdminAuth } from "@/lib/admin-auth";
import { AdminLayout } from "@/components/layout/admin-layout";
import { api } from "@/lib/api";
import type { Survey, SurveyQuestion, SurveyResult } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle,
} from "@/components/ui/dialog";
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from "@/components/ui/table";
import { Plus, Trash2, Eye } from "lucide-react";
import Link from "next/link";

export default function AdminSurveyPage() {
  const router = useRouter();
  const params = useParams();
  const moduleId = Number(params.moduleId);
  const { token, isLoading } = useAdminAuth();
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [results, setResults] = useState<SurveyResult[]>([]);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [resultDialogOpen, setResultDialogOpen] = useState(false);
  const [selectedResult, setSelectedResult] = useState<SurveyResult | null>(null);
  const [newQuestionText, setNewQuestionText] = useState("");
  const [newQuestionAnswerType, setNewQuestionAnswerType] = useState("text");
  const [newQuestionOptions, setNewQuestionOptions] = useState("");
  const [surveyTitle, setSurveyTitle] = useState("");

  useEffect(() => {
    if (!isLoading && !token) router.push("/admin/login");
  }, [token, isLoading, router]);

  useEffect(() => {
    if (token && moduleId) {
      api.get<Survey>(`/admins/modules/${moduleId}/survey`).then((res) => {
        setSurvey(res);
        if (res && res.title) setSurveyTitle(res.title);
      }).catch(() => setSurvey(null));
      api.get<{ data: SurveyResult[] }>(`/admins/modules/${moduleId}/survey/results`).then((res) => {
        setResults(res.data || []);
      }).catch(() => setResults([]));
    }
  }, [token, moduleId]);

  const handleCreateOrUpdate = async () => {
    if (!newQuestionText.trim()) {
      alert("请输入问题内容");
      return;
    }
    
    const questions = survey?.questions || [];
    const answerOptions = newQuestionOptions ? newQuestionOptions.split(",").map(s => s.trim()) : [];
    const newQuestions = [...questions, { 
      question: newQuestionText, 
      answerType: newQuestionAnswerType,
      answerOptions 
    }];
    
    try {
      await api.post(`/admins/modules/${moduleId}/survey`, { 
        title: surveyTitle || "问卷",
        required: true,
        questions: newQuestions 
      });
      setNewQuestionText("");
      setNewQuestionOptions("");
      setDialogOpen(false);
      const updated = await api.get<Survey>(`/admins/modules/${moduleId}/survey`);
      setSurvey(updated);
    } catch (error) {
      alert("添加问题失败: " + (error as Error).message);
    }
  };

  const handleDeleteSurvey = async () => {
    if (!confirm("确定删除该问卷？")) return;
    await api.delete(`/admins/modules/${moduleId}/survey`);
    setSurvey(null);
    setResults([]);
  };

  const viewResult = async (result: SurveyResult) => {
    setSelectedResult(result);
    setResultDialogOpen(true);
  };

  if (isLoading) return <AdminLayout><div className="flex items-center justify-center h-64">加载中...</div></AdminLayout>;

  return (
    <AdminLayout>
      <div className="space-y-6">
        <div>
          <Link href={`/admin/courses`} className="text-sm text-muted-foreground hover:underline">← 返回课程</Link>
          <h1 className="text-3xl font-bold mt-2">问卷管理</h1>
          <p className="text-muted-foreground">模块 #{moduleId}</p>
        </div>

        <div className="flex justify-between items-center">
          <h2 className="text-xl font-semibold">问卷内容</h2>
          <div className="flex gap-2">
            <Button onClick={() => setDialogOpen(true)}><Plus className="h-4 w-4 mr-2" />添加问题</Button>
            {survey && <Button variant="destructive" onClick={handleDeleteSurvey}>删除问卷</Button>}
          </div>
        </div>

        {survey && survey.questions.length > 0 ? (
          <Card>
            <CardHeader><CardTitle>问题列表</CardTitle></CardHeader>
            <CardContent>
              <div className="space-y-3">
                {survey.questions.map((q, i) => (
                  <div key={q.id} className="flex items-center justify-between p-3 border rounded-md">
                    <div>
                      <span className="text-sm text-muted-foreground mr-2">{i + 1}.</span>
                      <span>{q.question}</span>
                      {q.answerType && <span className="ml-2 text-xs text-blue-500">({q.answerType})</span>}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        ) : (
          <div className="text-center py-8 text-muted-foreground">暂无问卷</div>
        )}

        <h2 className="text-xl font-semibold">提交结果 ({results.length})</h2>
        <div className="bg-white rounded-lg border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>学生</TableHead>
                <TableHead>邮箱</TableHead>
                <TableHead className="w-[80px]">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {results.map((result) => (
                <TableRow key={result.id}>
                  <TableCell>{result.student?.name || `学生 #${result.studentId}`}</TableCell>
                  <TableCell>{result.student?.email || "-"}</TableCell>
                  <TableCell>
                    <Button variant="ghost" size="icon" onClick={() => viewResult(result)}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
              {results.length === 0 && (
                <TableRow><TableCell colSpan={3} className="text-center text-muted-foreground py-8">暂无提交</TableCell></TableRow>
              )}
            </TableBody>
          </Table>
        </div>

        {/* Add Question Dialog */}
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>添加问题</DialogTitle></DialogHeader>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label>问卷标题</Label>
                <Input value={surveyTitle} onChange={(e) => setSurveyTitle(e.target.value)} placeholder="输入问卷标题" />
              </div>
              <div className="space-y-2">
                <Label>问题内容</Label>
                <Input value={newQuestionText} onChange={(e) => setNewQuestionText(e.target.value)} placeholder="输入问题" />
              </div>
              <div className="space-y-2">
                <Label>问题类型</Label>
                <select 
                  value={newQuestionAnswerType} 
                  onChange={(e) => setNewQuestionAnswerType(e.target.value)} 
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                >
                  <option value="text">文本输入</option>
                  <option value="select">单选</option>
                  <option value="multiselect">多选</option>
                </select>
              </div>
              {(newQuestionAnswerType === "select" || newQuestionAnswerType === "multiselect") && (
                <div className="space-y-2">
                  <Label>选项（用逗号分隔）</Label>
                  <Input value={newQuestionOptions} onChange={(e) => setNewQuestionOptions(e.target.value)} placeholder="选项1, 选项2, 选项3" />
                </div>
              )}
              <Button onClick={handleCreateOrUpdate} className="w-full">添加</Button>
            </div>
          </DialogContent>
        </Dialog>

        {/* View Result Dialog */}
        <Dialog open={resultDialogOpen} onOpenChange={setResultDialogOpen}>
          <DialogContent>
            <DialogHeader><DialogTitle>问卷结果</DialogTitle></DialogHeader>
            {selectedResult && (
              <div className="space-y-3">
                {selectedResult.answers.map((a, i) => (
                  <div key={i} className="p-3 border rounded-md">
                    <div className="text-sm text-muted-foreground">问题 #{a.questionId}</div>
                    <div className="font-medium mt-1">{a.answer}</div>
                  </div>
                ))}
              </div>
            )}
          </DialogContent>
        </Dialog>
      </div>
    </AdminLayout>
  );
}
