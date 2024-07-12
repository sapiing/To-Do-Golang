"use client";

import { useState, useEffect } from "react";

export default function FormTask({ closeModal, refetchTasks, editingTask }) {
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    completed: false,
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (editingTask) {
      setFormData({
        title: editingTask.title,
        description: editingTask.description,
        completed: editingTask.completed,
      });
    }
  }, [editingTask]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      const token = localStorage.getItem("token");
      const method = editingTask ? "PUT" : "POST"; // Determine method based on editingTask
      const url = editingTask
        ? `http://localhost:8080/api/task/${editingTask.id}`
        : "http://localhost:8080/api/task";

      const response = await fetch(url, {
        method: method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...(editingTask && { id: editingTask.id }), // Include id only if editingTask exists
          title: formData.title,
          description: formData.description,
          completed: formData.completed,
        }),
      });

      if (!response.ok) {
        throw new Error("Network response was not ok.");
      }

      // Handle successful submission (e.g., clear form, show success message)
      setFormData({ title: "", description: "", completed: false});
      console.log(
        editingTask
          ? "Task updated successfully!"
          : "Task created successfully!"
      );

      closeModal(); //Call closeModal function after successful submission
      refetchTasks(); // Fetch tasks again to update the table.
    } catch (error) {
      setError(error.message);
      console.log(error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-md mx-auto mt-8">
      <div className="mb-4">
        <label
          htmlFor="title"
          className="block text-gray-700 text-sm font-bold mb-2"
        >
          Title
        </label>
        <input
          type="text"
          id="title"
          name="title"
          value={formData.title}
          onChange={handleChange}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>

      <div className="mb-6">
        <label
          htmlFor="description"
          className="block text-gray-700 text-sm font-bold mb-2"
        >
          Description
        </label>
        <textarea
          id="description"
          name="description"
          value={formData.description}
          onChange={handleChange}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline h-40"
        />
      </div>

      <div className="flex items-center justify-center">
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        >
          Submit
        </button>
      </div>
    </form>
  );
}
