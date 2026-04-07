"use client"
import { useMutation } from "@tanstack/react-query"
import { toast } from "sonner" // Đã cài sẵn trong project template

const createPost = async (newPost: any) => {
  const res = await fetch("https://jsonplaceholder.typicode.com/posts", {
    method: "POST",
    body: JSON.stringify(newPost),
    headers: { "Content-type": "application/json; charset=UTF-8" },
  })
  return res.json()
}

export function CreatePostButton() {
  const { mutate, isPending } = useMutation({
    mutationFn: createPost,
    onSuccess: (data) => {
      toast.success("Tạo Post thành công ID: " + data.id)
    },
    onError: () => {
      toast.error("Thất bại rồi!")
    }
  })

  return (
    <button 
      disabled={isPending}
      onClick={() => mutate({ title: "Foo", body: "Bar", userId: 1 })}
      className="bg-blue-500 text-white px-4 py-2 rounded disabled:opacity-50"
    >
      {isPending ? "Đang gửi..." : "Thêm Post"}
    </button>
  )
}